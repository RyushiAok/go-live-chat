// type MessageType = "message" | "join_new_member"

import { typeGuardOf } from "@/utils/typeGuard"

// export enum MessageType {
//     message = "message",
//     join_new_member = "join_new_member" 
// }  

export interface IReceiveMessage {
    type :  "message" ,
    data : {
        message : string
    }
}
export interface IJoinNewMember {
    type : "join_new_member",
    data : {
        name : string
    }
}
export interface ISendMessage {
    type : "send",
    data : {
        message : string
    }
}

 