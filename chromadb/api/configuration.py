from abc import abstractmethod
import json
from overrides import override
from typing import (
    Any,
    ClassVar,
    Dict,
    List,
    Optional,
    Protocol,
    Union,
    TypeVar,
)
from typing_extensions import Self
from multiprocessing import cpu_count

from chromadb.serde import JSONSerializable

# TODO: move out of API


class StaticParameterError(Exception):
    """Represents an error that occurs when a static parameter is set."""

    pass


ParameterValue = Union[str, int, float, bool, "Configuration"]


class ParameterValidator(Protocol):
    """Represents an abstract parameter validator."""

    @abstractmethod
    def __call__(self, value: ParameterValue) -> bool:
        """Returns whether the given value is valid."""
        raise NotImplementedError()


class ConfigurationDefinition:
    """Represents the definition of a configuration."""

    name: str
    validator: ParameterValidator
    is_static: bool
    default_value: ParameterValue

    def __init__(
        self,
        name: str,
        validator: ParameterValidator,
        is_static: bool,
        default_value: ParameterValue,
    ):
        self.name = name
        self.validator = validator
        self.is_static = is_static
        self.default_value = default_value


class ConfigurationParameter:
    """Represents a parameter of a configuration."""

    name: str
    value: ParameterValue

    def __init__(self, name: str, value: ParameterValue):
        self.name = name
        self.value = value

    def __repr__(self) -> str:
        return f"ConfigurationParameter({self.name}, {self.value})"

    def __eq__(self, __value: object) -> bool:
        if not isinstance(__value, ConfigurationParameter):
            return NotImplemented
        return self.name == __value.name and self.value == __value.value


T = TypeVar("T", bound="Configuration")


class Configuration(JSONSerializable["Configuration"]):
    """Represents an abstract configuration."""

    # The internal data structure used to store the parameters
    # All expected parameters must be present with defaults or None values at initialization
    parameter_map: Dict[str, ConfigurationParameter]
    definitions: ClassVar[Dict[str, ConfigurationDefinition]]

    def __init__(self, parameters: Optional[List[ConfigurationParameter]] = None):
        """Initializes a new instance of the Configuration class. Respecting defaults and
        validators."""
        self.parameter_map = {}
        if parameters is not None:
            for parameter in parameters:
                if parameter.name not in self.definitions:
                    raise ValueError(f"Invalid parameter name: {parameter.name}")

                definition = self.definitions[parameter.name]
                # Handle the case where we have a recursive configuration definition
                if isinstance(parameter.value, dict):
                    child_type = globals().get(parameter.value.get("_type", None))
                    if child_type is None:
                        raise ValueError(
                            f"Invalid configuration type: {parameter.value}"
                        )
                    parameter.value = child_type.from_json(parameter.value)
                if not isinstance(parameter.value, type(definition.default_value)):
                    raise ValueError(f"Invalid parameter value: {parameter.value}")

                validator = definition.validator
                if not validator(parameter.value):
                    raise ValueError(f"Invalid parameter value: {parameter.value}")
                self.parameter_map[parameter.name] = parameter
        # Apply the defaults for any missing parameters
        for name, definition in self.definitions.items():
            if name not in self.parameter_map:
                self.parameter_map[name] = ConfigurationParameter(
                    name=name, value=definition.default_value
                )

    def __repr__(self) -> str:
        return f"Configuration({self.parameter_map.values()})"

    def __eq__(self, __value: object) -> bool:
        if not isinstance(__value, Configuration):
            return NotImplemented
        return self.parameter_map == __value.parameter_map

    def get_parameters(self) -> List[ConfigurationParameter]:
        """Returns the parameters of the configuration."""
        return list(self.parameter_map.values())

    def get_parameter(self, name: str) -> Optional[ConfigurationParameter]:
        """Returns the parameter with the given name or None if it does not exist."""
        return self.parameter_map.get(name, None)

    def set_parameter(self, name: str, value: Union[str, int, float, bool]) -> None:
        """Sets the parameter with the given name to the given value."""
        if name not in self.definitions:
            raise ValueError(f"Invalid parameter name: {name}")
        definition = self.definitions[name]
        parameter = self.parameter_map[name]
        if definition.is_static:
            raise StaticParameterError(f"Cannot set static parameter: {name}")
        if not definition.validator(value):
            raise ValueError(f"Invalid value for parameter {name}: {value}")
        parameter.value = value

    @override
    def to_json_str(self) -> str:
        """Returns the JSON representation of the configuration."""
        return json.dumps(self.to_json())

    @classmethod
    @override
    def from_json_str(cls, json_str: str) -> Self:
        """Returns a configuration from the given JSON string."""
        try:
            config_json = json.loads(json_str)
        except json.JSONDecodeError:
            raise ValueError(
                f"Unable to decode configuration from JSON string: {json_str}"
            )
        return cls.from_json(config_json)

    @override
    def to_json(self) -> Dict[str, Any]:
        """Returns the JSON compatible dictionary representation of the configuration."""
        json_dict = {
            name: parameter.value.to_json()
            if isinstance(parameter.value, Configuration)
            else parameter.value
            for name, parameter in self.parameter_map.items()
        }
        # What kind of configuration is this?
        json_dict["_type"] = self.__class__.__name__
        return json_dict

    @classmethod
    @override
    def from_json(cls, json_map: Dict[str, Any]) -> Self:
        """Returns a configuration from the given JSON string."""
        if cls.__name__ != json_map.get("_type", None):
            raise ValueError(
                f"Trying to instantiate configuration of type {cls.__name__} from JSON with type {json_map['_type']}"
            )
        parameters = []
        for name, value in json_map.items():
            # Type value is only for storage
            if name == "_type":
                continue
            parameters.append(ConfigurationParameter(name=name, value=value))
        return cls(parameters=parameters)


class HNSWConfiguration(Configuration):
    definitions = {
        "space": ConfigurationDefinition(
            name="space",
            validator=lambda value: isinstance(value, str)
            and value in ["l2", "ip", "cosine"],
            is_static=True,
            default_value="l2",
        ),
        "ef_construction": ConfigurationDefinition(
            name="ef_construction",
            validator=lambda value: isinstance(value, int) and value >= 1,
            is_static=True,
            default_value=100,
        ),
        "ef_search": ConfigurationDefinition(
            name="ef_search",
            validator=lambda value: isinstance(value, int) and value >= 1,
            is_static=False,
            default_value=10,
        ),
        "num_threads": ConfigurationDefinition(
            name="num_threads",
            validator=lambda value: isinstance(value, int) and value >= 1,
            is_static=False,
            default_value=cpu_count(),  # By default use all cores available
        ),
        "M": ConfigurationDefinition(
            name="M",
            validator=lambda value: isinstance(value, int) and value >= 1,
            is_static=True,
            default_value=16,
        ),
        "resize_factor": ConfigurationDefinition(
            name="resize_factor",
            validator=lambda value: isinstance(value, float) and value >= 1,
            is_static=True,
            default_value=1.2,
        ),
        "batch_size": ConfigurationDefinition(
            name="batch_size",
            validator=lambda value: isinstance(value, int) and value >= 1,
            is_static=True,
            default_value=1000,
        ),
        "sync_threshold": ConfigurationDefinition(
            name="sync_threshold",
            validator=lambda value: isinstance(value, int) and value >= 1,
            is_static=True,
            default_value=100,
        ),
    }

    @classmethod
    def from_legacy_params(cls, params: Dict[str, Any]) -> Self:
        """Returns an HNSWConfiguration from a metadata dict containing legacy HNSW parameters. Used for migration."""

        # We maintain this map to avoid a circular import with HnswParams, and because then names won't change
        old_to_new = {
            "hnsw:space": "space",
            "hnsw:construction_ef": "ef_construction",
            "hnsw:search_ef": "ef_search",
            "hnsw:M": "M",
            "hnsw:num_threads": "num_threads",
            "hnsw:resize_factor": "resize_factor",
            "hnsw:batch_size": "batch_size",
            "hnsw:sync_threshold": "sync_threshold",
        }

        parameters = []
        for name, value in params.items():
            if name not in old_to_new:
                raise ValueError(f"Invalid legacy HNSW parameter name: {name}")
            new_name = old_to_new[name]
            parameters.append(ConfigurationParameter(name=new_name, value=value))
        return cls(parameters=parameters)


class CollectionConfiguration(Configuration):
    """The configuration for a collection."""

    definitions = {
        "hnsw_configuration": ConfigurationDefinition(
            name="hnsw_configuration",
            validator=lambda value: isinstance(value, HNSWConfiguration),
            is_static=False,
            default_value=HNSWConfiguration(),
        ),
    }
