import { Flex } from '@chakra-ui/react';
import { Link, useParams } from 'react-router-dom';

export default function Models() {
    let params = useParams();
    console.log('params', params)
    
    return (
      <Flex mt="48px">
        Models
      </Flex>
    )
}