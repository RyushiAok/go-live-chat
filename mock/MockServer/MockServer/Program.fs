 

[<EntryPoint>]
let main _ =
    //task{  MockWS.run() } |> ignore// websocketの接続確認
    async { MockWS.run() } |> Async.Start
    MockAPI.run() 
    0