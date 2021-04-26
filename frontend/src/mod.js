import React from "react"
import { ContentHandler, ControlBtn, CustomImage, DownloadableLink, fetchDataWithToken, useResponse } from "./base"

export function ModList(props) {
    const [content, ...state] = useResponse("/api/mods")

	return (
		<div className="mod-list">
 	       <ContentHandler contentState={state}>
				{
					content.model && content.model.map((value, index) => {
						return <Mod value={value} key={value.id || index} serverUptime={props.serverUptime}/>
					})
				}
			</ContentHandler>
		</div>
	)
}

export function Mod(props) {
	const path = "/api/mods/" + props.value.id
	const fileName = props.value.name
	const mod_time = new Date(props.value.mod_time)
	let className = props.className ? "download-link " + props.className : "download-link"

	return (
		<div className={props.serverUptime.getTime() < mod_time.getTime() ? "new mod" : "mod"}>
			<CustomImage name="jar_icon.png" className="jar-icon"/>
			<ControlBtn Element={DownloadableLink} path={path} fileName={fileName} className={className}>
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
		if (confirm('Are you sure you want to delete the file?')) {
			fetchDataWithToken(path, {method: "POST"})
			window.location.reload()
		}
	}
	return (
		<a onClick={onClick} className={fullClassName} {...attr}>{children}</a>
	)
}