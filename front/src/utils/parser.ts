import { IJoinNewMember, IReceiveMessage } from '@/types/api/room'//'@/types/api/room'
import * as O from 'fp-ts/Option'
import { typeGuardOf } from "@/utils/typeGuard"

module RoomApiTypeGuards {
    // https://typescriptbook.jp/reference/object-oriented/interface/instanceof-and-interfaces
    // https://zenn.dev/karan_coron/articles/492a866adcd7aa
    // https://zenn.dev/yonajs/articles/dc330d8d81b60f#%E5%9E%8B%E3%82%AC%E3%83%BC%E3%83%89(type-guard)%E3%81%A8%E3%81%AF%EF%BC%9F
 
    const MessageType = {
        message: "message",
        join_new_member : "join_new_member" 
    }  
    
    export const isIReceiveMessage = (obj: any): obj is IReceiveMessage => {
        try {
            const a: { type: string, data: any } = obj  
            if (a.type === MessageType.message) {
                return true
            } else {
                return false
            }
        }
        catch(e) {
            return false
        }
    }

    export const isIJoinNewMember = (obj: any): obj is IJoinNewMember => {
        try {
            const a: { type: string, data: any } = obj 
            if (a.type === MessageType.join_new_member) {
                return true
            } else {
                return false
            }
        }
        catch(e) {
            return false
        }
    }

}

export module IReceiveMessageHelper {
    export const tryParse = (data: any): O.Option<IReceiveMessage> => { 
        try {  
            if(RoomApiTypeGuards.isIReceiveMessage(data)) { 
                return O.some(data)
            }
            else{
                return O.none
            }
        }
        catch(e) {
            return O.none 
        } 
    }
}

export module IJoinNewMemberHeper {
    export const tryParse = (data: any): O.Option<IJoinNewMember> => { 
        try {  
            if(RoomApiTypeGuards.isIJoinNewMember(data)) { 
                return O.some(data)
            }
            else{
                return O.none
            }
        }
        catch(e) {
            return O.none 
        } 
    } 
} 