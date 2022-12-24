import { AxiosClient, GETApi } from '@/libs/api/AxiosClient'

const axios = AxiosClient("http://localhost:8080/api/room")
const get = GETApi(axios)
export default () => {  
    const GetInfo = (roomId : string) => get<string>(`?room=${roomId}`) 
    return { GetInfo }

}