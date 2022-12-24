module MockAPI

open Falco
open Falco.Routing
open Falco.HostBuilder

module Room = 
    type CreateRoom = {
        room_name : string
        max_count : int
    }
    type CreateRoomResp = {
        room_id : string 
        room_name : string
        max_count : int
        status : string 
    }
    type AddNewMemberResp = {
        
        user_id : string 
        room_id : string 
        room_name : string 
         
    }
     
    let createRoomHandler : HttpHandler = 
        let handleOk (query :CreateRoom) : HttpHandler =  
            Response.ofJson {|
                data = {
                    room_id= "649e9506c8b069296e6f"
                    room_name= query.room_name
                    max_count= query.max_count
                    status= "created"
                }
            |}  
        Request.mapJson handleOk

    let addNewMember : HttpHandler =
        let option =  System.Text.Json.JsonSerializerOptions()
        fun ctx -> task { 
            let r = Request.getQuery ctx 
            let roomId = r.GetString "room_id" ""
            if roomId = "" then 
                return! Response.ofEmpty ctx
            else  
                let! j = Request.getJsonOptions<AddNewMemberResp> option ctx 
                return! Response.ofJson {|
                    data = {
                        user_id = "1121" 
                        room_id = roomId
                        room_name = "room room room" 
                    }
                |} ctx
        }   

         
open Microsoft.AspNetCore.Builder
open Microsoft.AspNetCore.Hosting 
open Microsoft.Extensions.DependencyInjection 
open Microsoft.AspNetCore.Cors.Infrastructure 
open Microsoft.Extensions.Logging  

// ------------
// Register services
// ------------
let configureServices (services : IServiceCollection) =
    services.AddAntiforgery()
            .AddFalco() |> ignore

// ---------------------------------
// Config and Main
// ---------------------------------
 
let configureWebhost (endpoints : HttpEndpoint list) (webhost : IWebHostBuilder) =
    let configureCors (corsBuilder: CorsPolicyBuilder) : unit =
        corsBuilder
            .AllowAnyOrigin()
            .AllowAnyHeader()
            .AllowAnyMethod()
        |> ignore
    let configureApp (endpoints : HttpEndpoint list) (app : IApplicationBuilder) = 
        app
            .UseCors(fun options -> configureCors options) 
            .UseFalco(endpoints)  
        |> ignore 

    webhost 
        .ConfigureServices(fun service -> 
            service
                .AddCors() // https://stackoverflow.com/questions/46874706/unable-to-resolve-service-for-type-microsoft-aspnetcore-cors-infrastructure-ico
                .AddLogging(fun config -> config.SetMinimumLevel(LogLevel.Information) |> ignore) 
           |> ignore
        ) 
        .Configure(configureApp endpoints)
let run () =   

    webHost [||] { 
        configure configureWebhost 
        endpoints [ 
            get "/" (Response.ofPlainText "Hello World")
            get "/questions" (Response.ofPlainText "Hello World")
            post "/room/new" (Room.createRoomHandler)
            post "/member/new" (Room.addNewMember)
        ]
    }