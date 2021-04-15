import React from "react"
import { ContentHandler, ControlBtn, CustomImage, DownloadableLink, useResponse } from "./base"
import { getUserFromCookie } from "./user"

export function ModList(props) {
    const [content, ...state] = useResponse("/api/mods")
	return (
        <ContentHandler contentState={state}>
			<div className="mod-list">
				{
					content.model && content.model.map((value, index) => {
						return <Mod value={value} key={value.id || index}/>
					})
				}
			</div>
		</ContentHandler>
	)
}

export function Mod(props) {
	const path = "/api/mods/" + props.value.id
	const fileName = props.value.name
	return (
		<div className="mod">
			<CustomImage name="jar_icon.png" className="jar-icon"/>
			<ControlBtn Element={DownloadableLink} path={path} fileName={fileName} className="download-link">
				{props.value.name}
			</ControlBtn>
			<DeleteBtn path={path}>Del</DeleteBtn>
		</div>
    )
}

function DeleteBtn(props) {
	const {className, path, children, ...attr} = props
    const fullClassName = "delete-btn btn " + (className || "")
	function onClick(event) {
		event.preventDefault()
		if (confirm('Вы уверены что хотите удалить файл?')) {
			fetch(path, {
				headers: {
					'Content-Type': 'application/json',
					'X-Session-Token': 'Bearer ' + getUserFromCookie().token,
				},
				method: "POST",
			})
			window.location.reload()
		}
	}
	return (
		<a onClick={onClick} className={fullClassName} {...attr}>{children}</a>
	)
}