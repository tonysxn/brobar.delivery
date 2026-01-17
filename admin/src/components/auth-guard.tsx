"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState, useRef } from "react";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { AuthData } from "@/types/auth";
import axios from "axios";
import { BACKEND_URL } from "@/constants";

export function AuthGuard({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const pathname = usePathname();
    const [loginData, setLoginData] = useLocalStorage<AuthData | null>("auth", null);
    const [checking, setChecking] = useState(true);
    const [hydrated, setHydrated] = useState(false);
    const isCheckingRef = useRef(false);

    useEffect(() => {
        setHydrated(true);
    }, []);

    useEffect(() => {
        if (!hydrated) return;
        if (isCheckingRef.current) return;

        const publicRoutes = ["/login", "/register"];
        if (publicRoutes.includes(pathname)) {
            setChecking(false);
            return;
        }

        const checkAuth = async () => {
            isCheckingRef.current = true;

            if (!loginData) {
                router.replace("/login");
                return;
            }

            const accessExpiresAt = new Date(loginData.access.expires_at);
            const now = new Date();

            if (accessExpiresAt <= now) {
                setLoginData(null);
                router.replace("/login");
                return;
            }

            const timeLeftMs = accessExpiresAt.getTime() - now.getTime();
            const refreshIfLess = 12 * 60 * 60 * 1000;

            if (timeLeftMs <= refreshIfLess) {
                try {
                    const response = await axios.post(`${BACKEND_URL}/auth/refresh`, {
                        refresh_token: loginData.refresh.token,
                    });

                    if (response.status === 200) {
                        const data = response.data.data;
                        const newLoginData: AuthData = {
                            access: {
                                token: data.access.token,
                                expires_at: new Date(Date.now() + data.access.expires_in * 1000).toISOString(),
                            },
                            refresh: {
                                token: data.refresh.token,
                                expires_at: new Date(Date.now() + data.refresh.expires_in * 1000).toISOString(),
                            },
                            user: loginData.user,
                        };

                        setLoginData(newLoginData);
                    } else {
                        setLoginData(null);
                        router.replace("/login");
                        return;
                    }
                } catch (e) {
                    console.error(e);
                    setLoginData(null);
                    router.replace("/login");
                    return;
                }
            }

            // Check if user has admin role
            if (loginData.user.role_id !== "admin") {
                console.error("Access denied: User is not an admin");
                setLoginData(null);
                router.replace("/login");
                return;
            }

            isCheckingRef.current = false;
            setChecking(false);
        };

        checkAuth();
    }, [hydrated, pathname]); // Removed loginData from dependencies to prevent loops

    if (checking) return null;

    return <>{children}</>;
}
