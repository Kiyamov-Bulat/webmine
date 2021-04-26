import React from "react"
import ReactDOM from "react-dom"
import {User, saveUserInCookie, getUserFromCookie} from "./user"
import { LoginForm } from "./base";
import Home from "./home";

function App(props) {
    const [user, setUser] = React.useState(getUserFromCookie())
    const handleUserSet = (loggedUser) => {
        let newUser = new User(loggedUser)
        saveUserInCookie(newUser)
        setUser(newUser)
    }

    if (user.isLogged())
        return (<Home user={user}/>)
    else
        return (<LoginForm setUser={handleUserSet} />)
}

ReactDOM.render(
    <App/>,
    document.getElementById("root")
)