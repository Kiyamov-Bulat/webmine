import React from "react"
import ReactDOM from "react-dom"
import {User, saveUserInCookie, getUserFromCookie} from "./user"
import { LoginForm } from "./base";
import Home from "./home";

// const UserContext = React.createContext(null);

function App(props) {
    const [user, setUser] = React.useState(getUserFromCookie())
    React.useEffect(() => {
        user.isLogged() && saveUserInCookie(user)
    }, [user])
    const handleUserSet = React.useCallback((loggedUser) => {setUser(new User(loggedUser))}, [])
    if (user.isLogged())
        return (<Home user={user}/>)
    else
        return (<div><LoginForm setUser={handleUserSet} /></div>)
}

ReactDOM.render(
    <App/>,
    document.getElementById("root")
)