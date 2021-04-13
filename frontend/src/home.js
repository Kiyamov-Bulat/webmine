import React from "react"
import {ContentHandler, useResponse, List} from "./base"
import { ModList } from "./Mod"

export default function Home(props) {
    const response = useResponse("/api/mods")
    return (
        <ContentHandler contentState={response.slice(1)}>
            <ModList content={response[0].model}/>
        </ContentHandler>
    )
}