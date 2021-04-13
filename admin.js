import React from "react"
import {Switch, Route, NavLink, useRouteMatch} from "react-router-dom";
import {Section, LoginForm, ManagedForm} from "./base"
import { Article } from "./base";


export default function Admin(props) {
    const {url, path} = useRouteMatch()
    const {user, setUser, ...rest} = props
    if (user.isLogged()) {
        if (!user.isNormal()) {
            return (
                <div>
                    <div>
                        <ul>
                            <NavLinkLi to={`${url}/profile`}>Профиль</NavLinkLi>
                            <NavLinkLi to={`${url}/users`}>Пользователи</NavLinkLi>
                            <NavLinkLi to={`${url}/products`}>Продукты</NavLinkLi>
                        </ul>
                    </div>
                    <Switch>
                        <Route path={`${path}/profile`}>
                            <Profile user={user} setUser={setUser}/>
                        </Route>
                        <Route path={`${path}/users`}>
                            <Users/>
                        </Route>
                        <Route path={`${path}/products`}>
                            <Products/>
                        </Route>
                    </Switch>
                </div>
            )
        } else {
            return <div><h2>Эта секция недоступна обычным пользователям.</h2></div>
        }
    } else {
        return <div><LoginForm setUser={setUser}/></div>
    }
}

function Products(props) {
    return (
        <div>Not impelemented</div>
    )
}
function Users(props) {
    return (
        <div>Not impelemented</div>
    )
}
function Profile(props) {
    const [edit, setEdit] = React.useState(false)
    const user = props.user
    const {DeletedAt, token, ...userProps} = user
    const btnText = edit ? "назад" : "изменить"
    let element = null
    function handleClick(e) {
        e.preventDefault()
        setEdit(!edit)
    }
    if (edit) {
        element = <ProfileEditor user={userProps} setUser={props.setUser}/>
    } else { 
        element = createComponentsFromEntr(userProps, ([key, val]) => <div key={key}><Article title={key} text={val}/></div>)
    }
    return (
        <div>
            <button onClick={handleClick}>{btnText}</button>
            {element}
        </div>
    )
}

function ProfileEditor(props) {
    const {ID, ...user} = props.user
    function handleSubmit(e, formData) {
        const url = "/api/users/" + ID
        fetch(url, {
            body: JSON.stringify(formData),
            method: "PUT",
            headers: {
                'Content-Type': 'application/json',
                'X-Session-Token': `Bearer ${user.token}`
            }
        })
        .then(res => res.json())
        .then(
            (response) => {
                props.setUser(response.model) // in model field - new user, in user field - current user. It is so, because server handler handle every request with any user.
                window.location.reload()
            }, 
            (reason) => {
                console.log(reason)
            }
        )
    }
    return (
        <ManagedForm onSubmit={handleSubmit} content={user}>
            <input type="hidden" name="userID" value={ID}/>
            {createComponentsFromEntr(user, ([key, _]) =>  
                <input type="text" name={key} key={key}/>
            )}
            <input type="submit" value="Сохранить"/>
        </ManagedForm>
    )
}

function createComponentsFromEntr(object, func) {
    return Object.entries(object).map(arr => {
        if (typeof arr[1] != "function") {
            return func(arr)
        } 
    })
}

function NavLinkLi(props) {
    return (
        <li><NavLink to={props.to}>{props.children}</NavLink></li>
    )
}