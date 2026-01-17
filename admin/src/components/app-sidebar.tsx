"use client"

import * as React from "react"
import {
    IconChartBar,
    IconDashboard,
    IconDatabase,
    IconSettings,
    IconUsers,
    IconMeat,
    IconMessage,
    IconReceipt
} from "@tabler/icons-react"
import { BiBowlHot } from "react-icons/bi";

import { NavDocuments } from "@/components/nav-documents"
import { NavMain } from "@/components/nav-main"
import { NavSecondary } from "@/components/nav-secondary"
import { NavUser } from "@/components/nav-user"
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from "@/components/ui/sidebar"
import { AuthData } from "@/types/auth";
import { useLocalStorage } from "@/hooks/use-local-storage";

const data = {
    navMain: [
        {
            title: "Dashboard",
            url: "/dashboard",
            icon: IconDashboard,
        },
        {
            title: "Orders",
            url: "/dashboard/orders",
            icon: IconReceipt,
        },
        {
            title: "Analytics",
            url: "/dashboard/analytics",
            icon: IconChartBar,
        },
        {
            title: "Users",
            url: "/dashboard/users",
            icon: IconUsers,
        },
    ],
    navSecondary: [],
    subNav: [
        {
            name: "Categories",
            url: "/dashboard/categories",
            icon: IconDatabase,
        },
        {
            name: "Products",
            url: "/dashboard/products",
            icon: IconMeat,
        },
        {
            name: "Reviews",
            url: "/dashboard/reviews",
            icon: IconMessage,
        },
        {
            name: "Settings",
            url: "/dashboard/settings",
            icon: IconSettings,
        }
    ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
    const fallbackValue: AuthData | null = null;
    const [userData] = useLocalStorage<AuthData | null>(
        "auth",
        fallbackValue
    );

    return (
        <Sidebar collapsible="offcanvas" {...props}>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton
                            asChild
                            className="data-[slot=sidebar-menu-button]:!p-1.5"
                        >
                            <a href="https://brobar.delivery" target={"_blank"}>
                                <BiBowlHot className="!size-5" />
                                <span className="text-base font-semibold">brobar.delivery</span>
                            </a>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <NavMain items={data.navMain} />
                <NavDocuments items={data.subNav} />
                <NavSecondary items={data.navSecondary} className="mt-auto" />
            </SidebarContent>
            <SidebarFooter>
                <NavUser user={userData?.user} />
            </SidebarFooter>
        </Sidebar>
    )
}
