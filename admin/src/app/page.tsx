"use client"

import {AuthData} from "@/types/auth";
import {useLocalStorage} from "@/hooks/use-local-storage";
import {useEffect, useState} from "react";
import {useRouter} from "next/navigation";

export default function Home() {
    const router = useRouter()
    const fallbackValue: AuthData | null = null;
    const [loginData, setLoginDataReady] = useLocalStorage<AuthData | null>("auth", fallbackValue);
    const [hydrated, setHydrated] = useState(false)

    useEffect(() => {
        setHydrated(true)
    }, [])

    useEffect(() => {
        if (!hydrated) return
        if (loginData !== null) {
            router.push("/dashboard")
        } else {
            router.push("/login")
        }
    }, [loginData, hydrated, router]);

    return <div></div>
}
