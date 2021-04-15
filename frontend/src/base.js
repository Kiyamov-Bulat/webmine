import React from "react"
import { getUserFromCookie } from "./user"


export const IMAGE_STORAGE = "./static/images/"

export function ContentHandler(props) {
    const [isLoaded, error] = props.contentState
    if (error) {
        return <div>Произошла ошибка, пожалуйста, попробуйте перезагрузить страницу!{error.message}</div>
    } else if (!isLoaded) {
        return <div>Идет загрузка</div>
    } else { 
        return props.children
    }
}

export function ControlBtn(props) {
    const [isLoaded, setIsLoaded] =  React.useState(true)
    const [error, setError] = React.useState(null)
    const {path, Element, className, children, ...attr} = props
    let btnClassName = className ? className + " btn" : "btn"
    
    if (!isLoaded)
        btnClassName += btnClassName ? " " + "disabled" : "disabled"
    function fetchData(options) {
        const {method, headers, ...restOptions} = options || {}
        if (isLoaded)
        {
            setIsLoaded(false)
            return fetch(path, {
                method: method || "GET",
                headers: {
                    'Content-Type': 'application/json',
                    'X-Session-Token': 'Bearer ' + getUserFromCookie().token,
                    ...headers,
                },
                ...restOptions,
            })
        }
        return new Promise((resolve, reject) => reject('Wait another action!'))
    }
    function onLoadEnd() {
        setIsLoaded(true)
    }
    function onError(loadErr) {
        setError(loadErr)
        console.log(loadErr)
    }
    return (
        <Element onError={onError} onLoadEnd={onLoadEnd} fetchData={fetchData} className={btnClassName} {...attr}>
            {children}
        </Element>
    )
}

export function DownloadableLink(props) {
    const {fileName, fetchData, onLoadEnd, onError, children, ...attr} = props
	function downloadFile(event) {
		event.preventDefault()
		fetchData().then((response) => response.blob()
        ).then((blob) => {
			const url = window.URL.createObjectURL(
				new Blob([blob]),
			)
			const link = document.createElement('a')
			link.href = url
			link.setAttribute('download', fileName)
			document.body.appendChild(link)
			link.click()
			link.parentNode.removeChild(link)
            onLoadEnd()
		}, reason => onError(reason))
	}
	return (
		<a {...attr} onClick={downloadFile}>{children}</a>
    )
}

///////////////////////
// HOOKS
//////////////////////

export function useResponse(path, options) {
    const [content, setContent] = React.useState("")
    const [isLoaded, setIsLoaded] =  React.useState(false)
    const [error, setError] = React.useState(null)
    const {method, headers, ...restOptions} = options || {}

    React.useEffect(() => {
        const pathName = path
        fetch(pathName, {
            method: method || "GET",
            headers: {
                'Content-Type': 'application/json',
                'X-Session-Token': 'Bearer ' + getUserFromCookie().token,
                // ...headers,
            },
            ...restOptions,
        }
        ).then(res => res.json()
        ).then(
            response => {
                setContent(response)
                setIsLoaded(true)
            },
            reason => {
                setError(reason)
            }
        ).finally(() => {
            method == 'POST' && !error && window.location.reload()
        })
    }, [])
    return [content, isLoaded, error] 
}

// OTHER

export function CustomImage(props) {
    const {name, className, ...rest} = props
    return <img src={IMAGE_STORAGE + name} className={className} {...rest}/>
}

export function BackgroundImage(props) {
    const {className, style, _ref, ...rest} = props
    const newClassName = className && className + " image-div" || "image-div" 
    return <div ref={_ref} style={{"backgroundImage" : `url("${IMAGE_STORAGE}${props.name}")`, ...style}} className={newClassName} {...rest}></div>
}

export function LoginForm(props) {
    function handleSubmit(e, formData) {
        const url = "/api/login"
        fetch(url, {
            body: JSON.stringify(formData),
            method: "POST",
            headers: {'Content-Type': 'application/json'}
        }
        ).then(res => res.json()
        ).then(
            response => props.setUser(response.user), 
            error => {
                typeof props.handleError == "function" && props.handleError(error) || console.log(error)
            })
    }
    return (
        <ManagedForm className={"login-form"} onSubmit={handleSubmit} content={{email: "", password: ""}}>
            <label>Login<input type="text" name="email"/></label>
            <label>Password<input type="password" name="password"/></label>
            <input type="submit" value="Войти"/>
        </ManagedForm>
    )
}

export function ManagedForm(props) {
    const {content, onSubmit, children, ...rest} = props
    const [fieldsState, setFieldsState] = React.useState(content)
    function handleChange(e) {
        const name = e.target.name, value = e.target.value
        setFieldsState({...fieldsState, [name]:value})
    }
    function handleSubmit(e) {
        e.preventDefault()
        onSubmit(e, fieldsState)
    }
    return (
        <form {...rest} onSubmit={handleSubmit}>
            <ManagedFieldSet onChange={handleChange} fieldsState={fieldsState}>
                {children}
            </ManagedFieldSet>
        </form>
    )
}


export function ManagedFieldSet(props) {
    const {children, fieldsState, ...rest} = props
    const modChildren = React.Children.map(children, child => {
        if (child) {
            if (child.props.type != "submit") {
                return React.cloneElement(child, {
                    value: fieldsState[child.props.name],
                    ...rest
                })
            } else {
                return child
            }
        }
    })
    return (
        <fieldset>{modChildren}</fieldset>
    )
}

export function usePreloadImages(objects) {
    const preloadedImages = React.useRef()
    preloadedImages.current = []
    React.useEffect(() => objects.forEach((val) => {
        const pImage = new Image()
        pImage.src = IMAGE_STORAGE + val.imageName
        preloadedImages.current.push(pImage)
    }), [])
}