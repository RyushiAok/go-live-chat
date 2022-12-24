import * as O from 'fp-ts/Option'
import * as A from 'fp-ts/Array'
import { pipe } from 'fp-ts/lib/function'

import { useEffect } from 'react'
import { useWebSocket } from 'react-use-websocket/dist/lib/use-websocket'
import { ReadyState } from 'react-use-websocket'

export default (
    url: string,
    handlers: ((string: string) => O.Option<void>)[],
    handleUnexeptedMessage: (message: string) => void
) => {
    const { sendMessage, lastJsonMessage, readyState } = useWebSocket(url) 

    useEffect(() => {
        if (lastJsonMessage !== null) {
            const json: any = lastJsonMessage
            pipe(
                handlers,
                A.findFirstMap((f) => f(json)),
                O.match(
                    () => { handleUnexeptedMessage(json) },
                    (_) => { },
                )
            )
        }
    }, [lastJsonMessage])

    const send = (data: string) => sendMessage(data)

    const connectionStatus = {
        [ReadyState.CONNECTING]: "Connecting",
        [ReadyState.OPEN]: "Open",
        [ReadyState.CLOSING]: "Closing",
        [ReadyState.CLOSED]: "Closed",
        [ReadyState.UNINSTANTIATED]: "Uninstantiated",
    }[readyState];
 
    return { send, connectionStatus }
}