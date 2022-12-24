import axios, {AxiosInstance} from 'axios'
import * as TE from 'fp-ts/TaskEither'

export const AxiosClient = (url : string): AxiosInstance => {
  return axios.create({
    baseURL: url, 
    headers: {
        'Content-Type':  'application/json'
    } 
  })
}

export const POSTApi =  (client: AxiosInstance) =>  <In, Out>(path: string) => (request: In) : TE.TaskEither<Error, Out> => {
    return TE.tryCatch(
        async () => { 
            const { data } = await client.post(path, JSON.stringify(request))
            return data.data
        },
        (reason) => new Error(String(reason)),
    )
}

export const GETApi = (client: AxiosInstance) => <Out> (param: string) : TE.TaskEither<Error, Out> => {
    return TE.tryCatch(
        async () => { 
            const { data } = await client.get(param)
            return data.data
        },
        (reason) => new Error(String(reason)),
    )
}  
