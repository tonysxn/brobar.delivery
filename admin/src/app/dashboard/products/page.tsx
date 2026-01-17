"use client"

import {useEffect, useState} from "react"
import {QueryClient} from "@tanstack/react-query"
import {PersistQueryClientProvider} from "@tanstack/react-query-persist-client"
import dynamic from "next/dynamic"
import type {ComponentType} from "react"

const ProductsPage = dynamic(() =>
    import("./_products").then(mod => mod.default as ComponentType), {ssr: false}
)


const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 5 * 60 * 1000,
            refetchOnWindowFocus: false,
            refetchOnMount: false,
            refetchOnReconnect: false,
        },
    },
})

export default function ProductsPageWrapper() {
    const [persister, setPersister] = useState<any>(null)

    useEffect(() => {
        async function initPersister() {
            if (typeof window === "undefined") return
            const mod = await import("@tanstack/query-async-storage-persister")
            const _persister = mod.createAsyncStoragePersister({
                storage: window.localStorage,
            })
            setPersister(_persister)
        }

        initPersister()
    }, [])

    if (!persister) return null

    return (
        <PersistQueryClientProvider client={queryClient} persistOptions={{persister}}>
            <ProductsPage/>
        </PersistQueryClientProvider>
    )
}
