import React from "react"
import { DownloadableLink , ControlBtn, fetchDataWithToken, ProgressBar, recursiveMap, ToLoadContext} from "./base"
import { ModList } from "./mod"

export default function Home(props) {
    const [content, setContent] = React.useState({})
    const [isLoaded, setIsLoaded] =  React.useState(true)
    const [error, setError] = React.useState(null)
    const uptime = content.model && content.model.uptime && new Date(content.model.uptime) || new Date()

    React.useEffect(() => {
        fetchDataWithToken("/api/server"
        ).then(res => res.json()
        ).then(
            response => {
                setContent(response)
                setIsLoaded(true)
            },
            reason => setError(reason))
    }, [props.user])

    return (
        <div id="home">
            <LoadHandler>
                <ModList serverUptime={uptime} />
                <ControlPanel server={content.model || {}} updateServer={setContent}/>
            </LoadHandler>
        </div>
    )
}

function LoadHandler(props) {
    const [toLoad, setToLoad] = React.useState({})

    function handleLoad(response) {
        const contentLength = +response.headers.get('Content-Length')
        const reader = response.body.getReader()
        let receivedLength = 0
        let chunks = []
        let id = Date.now().toString()
        let loader = reader.read().then(function rdChunk({done, value}) {
            if (done) return
            chunks.push(value)
            receivedLength += value.length
            setToLoad({"id": id, "item": <ProgressBar id={id} size={contentLength} progress={receivedLength} key={id}/>})
            return reader.read().then(rdChunk)
        })
        setTimeout(() => setToLoad({"id": id, "item": null}), 2000)
        return loader.then(() => {
            return new Promise((resolve, reject) => { 
                if (receivedLength > 0)
                    resolve(new Blob(chunks))
                else 
                    reject('failed to fetch data!')
            })
        })
    }
    return (
        <ToLoadContext.Provider value={{toLoad, handleLoad}}>
            {props.children}
        </ToLoadContext.Provider>
    )
}

function StatusPanel(props) {
    const [bars, setBars] = React.useState({})
    const {toLoad} = React.useContext(ToLoadContext)

    React.useEffect(() => {
        if (toLoad.item == null) {
            let tmp = {...bars}
            delete tmp[toLoad.id]
            setBars(tmp);
        } else
            setBars({...bars, ...{[toLoad.id]: toLoad.item}})
    }, [toLoad])
    return (
        <div className="status-panel-wrap">
            <div className="status-panel">
                {Object.values(bars || {})}
            </div>
        </div>
    )
}

function ControlPanel(props) {
    return (
        <div className="control-panel">
            <ControlBtn Element={FileUploadBtn} path="/api/mods" name="file">MOD_UPLOAD</ControlBtn>
            <ControlBtn Element={DownloadableLink} path="/api/archive/mods" fileName="mods.zip" >MODS_ARCHIVE</ControlBtn>
            <ControlBtn Element={DownloadableLink} path="/api/archive/data" fileName="server_data.zip">DATA_ARCHIVE</ControlBtn>
            <ControlBtn Element={ReloadServerBtn} path="/api/server/reload" onLoadEnd={props.updateServer} >RELOAD_SERVER</ControlBtn>
            <Board server={props.server}>
                <StatusPanel/>
            </Board>
        </div>
    )
}

const STATE_RUNNING = "Running"
const STATE_RESTARTING = "Restarting"

function Board(props) {
    const server = props.server
    
    return (
        <div className="board-wrap">
            <div className="board">
                <UptimeTable uptime={server.uptime}/>
                <StateTable state={server.state}/>
            </div>
            {props.children}
        </div>
    )
}

function StateTable(props) {
    let state = props.state || ""

    switch (state) {
        case STATE_RUNNING:
            state += " ✔️"
            break
        case STATE_RESTARTING:
            state += " ⌛"
            break
        default:
            state += " ⚠️"
    }
    return (
        <div className="status">
            <h2>- Status -</h2>
            <p>{state}</p>
        </div>
    )
}

function UptimeTable(props) {
    const uptime = props.uptime ? new Date(props.uptime) : new Date()
    const uptimeStr = `${uptime.toDateString()}  ${uptime.getHours()}:${uptime.getMinutes()}:${uptime.getSeconds()}`
    
    return (
        <div className="uptime">
            <h2>- Uptime -</h2>
            <p>{uptimeStr}</p>
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
            headers: { "Content-Type": undefined },
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