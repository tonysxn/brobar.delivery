"use client"

import { useEffect, useState } from "react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import axios from "axios";
import { BACKEND_URL } from "@/constants";
import showErrorToast, { showSuccessToast } from "@/components/toast";
import { AuthData } from "@/types/auth";
import { useLocalStorage } from "@/hooks/use-local-storage";
import { useRouter } from "next/navigation";

export function LoginForm({
    className,
    ...props
}: React.ComponentProps<"div">) {
    const router = useRouter()

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [hydrated, setHydrated] = useState(false)

    const fallbackValue: AuthData | null = null;
    const [loginData, setLoginData] = useLocalStorage<AuthData | null>(
        "auth",
        fallbackValue
    );

    useEffect(() => {
        setHydrated(true)
    }, [])

    useEffect(() => {
        if (!hydrated) return
        if (loginData !== null) {
            router.push("/dashboard");
        }
    }, [loginData, hydrated, router]);

    const login = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()

        try {
            const response = await axios.post(BACKEND_URL + "/auth/login", {
                email: email,
                password: password,
            });

            if (response.status === 200) {
                const data = response.data.data;

                // Check if user has admin role
                if (data.user.role_id !== "admin") {
                    showErrorToast("Access denied. Admin privileges required.");
                    return;
                }

                showSuccessToast("Login successful!");

                const now = Date.now();

                const expires_at_access = new Date(now + data.access.expires_in * 1000).toISOString();
                const expires_at_refresh = new Date(now + data.refresh.expires_in * 1000).toISOString();

                const loginResponse: AuthData = {
                    access: {
                        token: data.access.token,
                        expires_at: expires_at_access,
                    },
                    refresh: {
                        token: data.refresh.token,
                        expires_at: expires_at_refresh,
                    },
                    user: data.user
                };

                setLoginData(loginResponse);
                router.push("/");
            }
        } catch (e) {
            if (e.response && e.response.status === 401) {
                showErrorToast("Invalid email or password")
            }
        }
    }

    return (
        <div className={cn("flex flex-col gap-6", className)} {...props}>
            <Card>
                <CardHeader className="text-center">
                    <CardTitle className="text-xl">Welcome back</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={login}>
                        <div className="grid gap-6">
                            <div className="grid gap-6">
                                <div className="grid gap-3">
                                    <Label htmlFor="email">Email</Label>
                                    <Input
                                        id="email"
                                        type="email"
                                        placeholder="m@example.com"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        required
                                    />
                                </div>
                                <div className="grid gap-3">
                                    <div className="flex items-center">
                                        <Label htmlFor="password">Password</Label>
                                    </div>
                                    <Input
                                        id="password"
                                        type="password"
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        required
                                    />
                                </div>
                                <Button type="submit" className="w-full">
                                    Login
                                </Button>
                            </div>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    )
}
