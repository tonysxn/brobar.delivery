"use client"

import * as React from "react"
import { IconDotsVertical, IconPlus } from "@tabler/icons-react"
import { ColumnDef } from "@tanstack/react-table"
import { useIsMobile } from "@/hooks/use-mobile"
import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
    DialogClose,
} from "@/components/ui/dialog"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import axios from "axios"
import { BACKEND_URL } from "@/constants"
import { DataTable } from "@/components/ui/datatable"
import { useLocalStorage } from "@/hooks/use-local-storage"
import { AuthData } from "@/types/auth"
import showErrorToast, { showSuccessToast } from "@/components/toast"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { SortableHeader } from "@/components/ui/sortable-header"
import { Category } from "@/app/dashboard/categories/types"
import { Pagination } from "@/app/dashboard/types"

export default function Page() {
    const [data, setData] = React.useState<Category[]>([])
    const [pagination, setPagination] = React.useState<Pagination>({
        page: 1,
        total: 0,
        limit: 10,
        orderBy: "sort",
        orderDir: "asc",
    })
    const [queryParams, setQueryParams] = React.useState({
        page: 1,
        limit: 10,
        orderBy: "sort",
        orderDir: "asc" as "asc" | "desc",
    })
    const [editingItem, setEditingItem] = React.useState<Category | null>(null)
    const [dialogOpen, setDialogOpen] = React.useState(false)
    const [formData, setFormData] = React.useState({
        name: "",
        slug: "",
        icon: "",
        sort: 0,
    })
    const [deleteId, setDeleteId] = React.useState<string | null>(null)
    const [alertOpen, setAlertOpen] = React.useState(false)
    const isMobile = useIsMobile()
    const [user] = useLocalStorage<AuthData | null>("auth", null)

    const fetchData = React.useCallback(async () => {
        try {
            const res = await axios.get(`${BACKEND_URL}/categories`, {
                params: {
                    page: queryParams.page,
                    limit: queryParams.limit,
                    order_by: queryParams.orderBy,
                    order_dir: queryParams.orderDir,
                },
            })
            setData(res.data.data)
            setPagination({
                page: res.data.pagination.page,
                total: res.data.pagination.total_count,
                limit: res.data.pagination.limit,
                orderBy: res.data.pagination.order_by,
                orderDir: res.data.pagination.order_dir,
            })
        } catch {
            showErrorToast("Failed to fetch categories")
        }
    }, [queryParams])

    React.useEffect(() => {
        fetchData()
    }, [fetchData])

    React.useEffect(() => {
        if (editingItem) {
            setFormData({
                name: editingItem.name,
                slug: editingItem.slug,
                icon: editingItem.icon,
                sort: editingItem.sort,
            })
        }
    }, [editingItem])

    function openCreateDialog() {
        setEditingItem(null)
        setFormData({ name: "", slug: "", icon: "", sort: 0 })
        setDialogOpen(true)
    }

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault()
        try {
            if (editingItem) {
                const res = await axios.put(
                    `${BACKEND_URL}/categories/${editingItem.id}`,
                    formData,
                    { headers: { Authorization: `Bearer ${user?.access.token}` } }
                )
                const updatedItem = res.data.data
                setData((prev) => prev.map((cat) => (cat.id === editingItem.id ? updatedItem : cat)))
                setEditingItem(updatedItem)
                showSuccessToast("Category updated successfully")
            } else {
                const res = await axios.post(
                    `${BACKEND_URL}/categories`,
                    formData,
                    { headers: { Authorization: `Bearer ${user?.access.token}` } }
                )
                const newItem = res.data.data
                setData((prev) => [...prev, newItem])
                setEditingItem(newItem)
                showSuccessToast("Category created successfully")
            }
            fetchData()
            setDialogOpen(false)
        } catch (error: any) {
            const errorMessage = error.response?.data?.error || (editingItem ? "Failed to update category" : "Failed to create category")
            showErrorToast(errorMessage)
        }
    }

    async function handleDelete() {
        if (!deleteId) return
        try {
            await axios.delete(`${BACKEND_URL}/categories/${deleteId}`, {
                headers: { Authorization: `Bearer ${user?.access.token}` },
            })
            setData((prev) => prev.filter((cat) => cat.id !== deleteId))
            setDeleteId(null)
            setAlertOpen(false)
            showSuccessToast("Category deleted successfully")
        } catch (error: any) {
            const errorMessage = error.response?.data?.error || "Failed to delete category"
            showErrorToast(errorMessage)
        }
    }

    function handleSortChange(orderBy: string, orderDir: "asc" | "desc") {
        setQueryParams((prev) => ({ ...prev, orderBy, orderDir, page: 1 }))
    }

    function handlePaginationChange(page: number, limit: number) {
        setQueryParams((prev) => ({ ...prev, page, limit }))
    }

    const columns: ColumnDef<Category>[] = [
        {
            accessorKey: "name",
            header: () => (
                <SortableHeader
                    columnId="name"
                    currentOrderBy={pagination.orderBy}
                    currentOrderDir={pagination.orderDir}
                    onSortChange={handleSortChange}
                >
                    Name
                </SortableHeader>
            ),
            enableHiding: false,
        },
        {
            accessorKey: "slug",
            header: () => (
                <SortableHeader
                    columnId="slug"
                    currentOrderBy={pagination.orderBy}
                    currentOrderDir={pagination.orderDir}
                    onSortChange={handleSortChange}
                >
                    Slug
                </SortableHeader>
            ),
            enableHiding: false,
        },
        {
            accessorKey: "sort",
            header: () => (
                <SortableHeader
                    columnId="sort"
                    currentOrderBy={pagination.orderBy}
                    currentOrderDir={pagination.orderDir}
                    onSortChange={handleSortChange}
                >
                    Sort
                </SortableHeader>
            ),
            enableHiding: false,
        },
        {
            id: "actions",
            cell: ({ row }) => (
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon"
                            className="data-[state=open]:bg-muted text-muted-foreground flex size-8">
                            <IconDotsVertical />
                            <span className="sr-only">Open menu</span>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-32">
                        <DropdownMenuItem
                            onClick={() => {
                                setEditingItem(row.original)
                                setDialogOpen(true)
                            }}
                        >
                            Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem
                            className="text-destructive"
                            onClick={() => {
                                setDeleteId(row.original.id)
                                setAlertOpen(true)
                            }}
                        >
                            Delete
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            ),
        },
    ]

    return (
        <>
            <DataTable
                data={data}
                columns={columns}
                pagination={pagination}
                customButtons={[
                    <Button
                        key="create-btn"
                        variant="outline"
                        size="sm"
                        className="flex items-center text-white border-white hover:bg-white/10"
                        onClick={openCreateDialog}
                    >
                        <IconPlus />
                        <span className="hidden lg:inline">Create</span>
                        <span className="lg:hidden leading-none">Create</span>
                    </Button>,
                ]}
                onPaginationChange={handlePaginationChange}
            />

            <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                <DialogContent
                    className={`h-full lg:h-auto w-full max-w-md p-6 shadow-lg overflow-auto ${isMobile ? "max-w-full" : ""}`}>
                    <DialogHeader>
                        <DialogTitle>{editingItem ? editingItem.name : "Create New Category"}</DialogTitle>
                        <DialogClose className="absolute top-4 right-4" />
                    </DialogHeader>
                    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
                        <div className="flex flex-col gap-3">
                            <Label>Name</Label>
                            <Input value={formData.name}
                                onChange={(e) => setFormData((prev) => ({ ...prev, name: e.target.value }))} />
                        </div>
                        <div className="flex flex-col gap-3">
                            <Label>Slug</Label>
                            <Input value={formData.slug}
                                onChange={(e) => setFormData((prev) => ({ ...prev, slug: e.target.value }))} />
                        </div>
                        <div className="flex flex-col gap-3">
                            <Label>Icon</Label>
                            <Input value={formData.icon}
                                onChange={(e) => setFormData((prev) => ({ ...prev, icon: e.target.value }))} />
                        </div>
                        <div className="flex flex-col gap-3">
                            <Label>Sort</Label>
                            <Input
                                type="number"
                                value={formData.sort}
                                onChange={(e) => setFormData((prev) => ({
                                    ...prev,
                                    sort: parseInt(e.target.value) || 0
                                }))}
                            />
                        </div>
                        <DialogFooter className="mt-4">
                            <Button className="w-full" type="submit">
                                Submit
                            </Button>
                        </DialogFooter>
                    </form>
                </DialogContent>
            </Dialog>

            <AlertDialog open={alertOpen} onOpenChange={setAlertOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete Category</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to delete this category? This action cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleDelete}
                            className="bg-destructive text-white hover:bg-destructive/90">
                            Delete
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}
