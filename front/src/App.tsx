import { useState } from 'react' 
import './App.css' 
import useGroupChat from './hooks/useGroupChat'

function App() {
  const [count, setCount] = useState(0)
  const {send, connectionStatus} = useGroupChat(
    'room_1234',
    (newMessage) => {

      console.log("IReceiveMessage: ", newMessage)
    
    },
    (newMessage) => {

      console.log("IJoinNewMember: ", newMessage)
    
    },
    (a) => console.log("invalid, ", a)
  ) 
  return (
    <div className="App"> 
      <p>{connectionStatus}</p>
    </div>
  )
}

export default App
