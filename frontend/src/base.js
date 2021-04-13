import React from "react"

export const IMAGE_STORAGE = "./static/images/"

export function Section(props) {
    const {className, title, children, _ref, ...rest} = props
    const header = title && <h1>{props.title}</h1>
    const sectionRef = React.useCallback((section) => {
        if (section !== null) {
            if (window.location.hash == "#" + className) section.scrollIntoView()
            if (typeof _ref == "function") _ref(section)
        }
    }, [])
    return (
        <section className={className + " section"} id={className} ref={sectionRef} {...rest}>
            {header}
            {children}
        </section>
    )
}

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

export function Logo(props) {
    return (
        <h1 className="link-logo"><a href="/">{props.children}</a></h1>
    )
}

export function List(props) { 
    const {className = "", items = [], specListItemComponent, ...rest} = props
    const ListUnit = specListItemComponent || ListItem
    return (
        <ul {...rest} className={className}>
            {props.items.map((value, index) => {
                return <ListUnit content={value} key={value.id || index}/>
            })}
        </ul>
    )
}

export function ListItem(props) {
    const {content, ...rest} = props
    return (
        <li {...rest}>
            <Article {...content}/>
        </li>
    )
}

export function LinkListItem(props) {
    const {content, ...rest} = props,
          {href, ...articleContent} = content
    return (
        <li {...rest}>
            <a href={href}>
                <Article {...articleContent}/>
            </a>
        </li>
    )
}

export function Article(props) {
    const {title, imageName, text} = props,
          aTitle = title && <h2>{title}</h2>,
          aImage = imageName && <BackgroundImage name={imageName}/>,
          aText  = text && <p>{text}</p>
    return (
        <React.Fragment>
            {aTitle}
            {aImage}
            {aText}
        </React.Fragment>
    )
}

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


export function DownloadableLink(props) {
    const {path, fileName, children, ...linkAttrs} = props
	function downloadFile(event) {
		event.preventDefault()
		fetch(path, {
		method: 'GET',
		headers: {
		'Content-Type': 'application/json',
		},
		})
		.then((response) => response.blob())
		.then((blob) => {
			// Create blob link to download
			const url = window.URL.createObjectURL(
				new Blob([blob]),
			)
			const link = document.createElement('a')
			link.href = url
			link.setAttribute('download', fileName)

			// Append to html link element page
			document.body.appendChild(link)

			// Start download
			link.click()

			// Clean up and remove the link
			link.parentNode.removeChild(link)
		})
	}
	return (
		<a {...linkAttrs} onClick={downloadFile}>{children}</a>
    )
}

///////////////////////
// HOOKS
//////////////////////

export function useResponse(path) {
    const [content, setContent] = React.useState("")
    const [isLoaded, setIsLoaded] =  React.useState(false)
    const [error, setError] = React.useState(null)

    React.useEffect(() => {
        const pathName = path
        fetch(pathName, {
            method: "GET"
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
        )
    }, [])
    return [content, isLoaded, error] 
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