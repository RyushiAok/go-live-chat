import * as O from 'fp-ts/Option' 
import { flow } from 'fp-ts/lib/function' 
import useWebSocket from './useWebSocket'
import { IJoinNewMember, IReceiveMessage } from '@/types/api/room' 
import { IJoinNewMemberHeper, IReceiveMessageHelper } from '@/utils/parser' 

// websocketの接続をもち，受信したメッセージをハンドラに投げるhook 
export default (  
    roomId : string,
    handleNewMessage: (newMessage: IReceiveMessage) => void,  
    handleJoinNewMember: (member: IJoinNewMember) => void,
    handleUnexeptedMessage : (message: string) => void
) => {
    const handlers: ((string: string) => O.Option<void>)[] = [
        flow(IReceiveMessageHelper.tryParse, O.map(handleNewMessage)),
        flow(IJoinNewMemberHeper.tryParse ,  O.map(handleJoinNewMember))
    ]

    
    const {send, connectionStatus} = useWebSocket(
        `ws://localhost:8080/ws/${roomId}`,
        handlers,
        handleUnexeptedMessage
    ) 
    return { send, connectionStatus } 
}

// パターンマッチの失敗の責任を分離（型ガードを書くのがだるい）
// コンポーネント側で色々書いても良し，フックに全部書いても良し
// Option型で合成性を向上
// 型以外の条件も増やせる