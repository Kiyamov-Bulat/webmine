import React from "react"
import { DownloadableLink, useResponse } from "./base"

export function ModList(props) {
	return (
		<div className="mod-list">
			{props.content.map((value, index) => {
				return <Mod value={value} key={value.id || index}/>
			})}
		</div>
	)
}

export function Mod(props) {
	const path = "/api/mods/" + props.value.id
	const fileName = props.value.name
	return (
		<DownloadableLink href="" path={path} fileName={fileName} className="mod">{props.value.name}</DownloadableLink>
    )
}