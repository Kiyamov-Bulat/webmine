import React from "react"
import { DownloadableLink , ControlBtn} from "./base"
import { ModList } from "./mod"
import { getUserFromCookie } from "./user"


export default function Home(props) {
    return (
        <div>
            <ModList/>
            <ControlPanel/>
        </div>
    )
}

function ControlPanel(props) {
    return (
        <div className="control-panel">
            <ControlBtn Element={FileUploadBtn} path="/api/mods" name="file">MOD_UPLOAD</ControlBtn>
            <ControlBtn Element={DownloadableLink} path="/api/archive/mods" fileName="mods.zip" >MODS_ARCHIVE</ControlBtn>
            <ControlBtn Element={DownloadableLink} path="/api/archive/data" fileName="server_data.zip">DATA_ARCHIVE</ControlBtn>
            <ControlBtn Element={ReloadServerBtn} path="/api/server/reload">RELOAD_SERVER</ControlBtn>
        </div>
    )
}

function FileUploadBtn(props) {
    const {name, fetchData, onLoadEnd, onError, children, ...attr} = props

    function onChange(event) {
        event.preventDefault()
        const formData = new FormData()
        const file = event.target.files[0]
        formData.append('file', file)
        fetchData({
            method: "POST",
            body: formData,
        }
        ).then(res => res.json()
        ).then(
            response => onLoadEnd(response),
            reason   => onError(reason)
        ).finally(() => {
            window.location.reload()
        })
    }
    return (
        <label {...attr}>
            <input onChange={onChange} type="file" name={name}/>
            {children}
        </label>
    )
}

function ReloadServerBtn(props) {
    const {fetchData, onLoadEnd, onError, children, ...attr} = props
    function reload(event) {
		event.preventDefault()
        fetchData().then(res => res.json()
        ).then(
            response => onLoadEnd(response),
            reason   => onError(reason)
        )
    }
    return (
        <a onClick={reload} href="" {...attr}>{children}</a>
    )
}