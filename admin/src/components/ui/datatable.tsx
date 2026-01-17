import * as React from "react"
import {
    closestCenter,
    DndContext,
    type DragEndEvent,
    KeyboardSensor,
    MouseSensor,
    TouchSensor,
    type UniqueIdentifier,
    useSensor,
    useSensors,
} from "@dnd-kit/core"
import {
    ColumnDef,
    ColumnFiltersState,
    flexRender,
    getCoreRowModel,
    getFacetedRowModel,
    getFacetedUniqueValues,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    SortingState,
    useReactTable,
    VisibilityState,
} from "@tanstack/react-table"
import {arrayMove, SortableContext, verticalListSortingStrategy} from "@dnd-kit/sortable"
import {Tabs, TabsContent} from "@/components/ui/tabs"
import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import {Button} from "@/components/ui/button"
import {
    IconChevronDown,
    IconChevronLeft,
    IconChevronRight,
    IconChevronsLeft,
    IconChevronsRight,
    IconLayoutColumns,
} from "@tabler/icons-react"
import {restrictToVerticalAxis} from "@dnd-kit/modifiers"
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table"
import {Label} from "@/components/ui/label"
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select"

function useTableState<T>(data: T[]) {
    const [sorting, setSorting] = React.useState<SortingState>([])
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([])
    const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({})
    const [rowSelection, setRowSelection] = React.useState({})
    const [pagination, setPagination] = React.useState({pageIndex: 0, pageSize: 10})

    return {
        sorting,
        setSorting,
        columnFilters,
        setColumnFilters,
        columnVisibility,
        setColumnVisibility,
        rowSelection,
        setRowSelection,
        pagination,
        setPagination,
        tableState: {
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
            pagination,
        },
    }
}

export function DataTable<TData extends { id: string }>({
                                                            data: initialData,
                                                            columns,
                                                            customButtons,
                                                            pagination,
                                                            onPaginationChange,
                                                        }: {
    data: TData[]
    columns: ColumnDef<TData>[]
    customButtons?: React.ReactNode[]
    pagination?: { page: number; limit: number; total: number }
    onPaginationChange?: (page: number, limit: number) => void
}) {
    const [data, setData] = React.useState(initialData)

    React.useEffect(() => {
        setData(initialData)
    }, [initialData])

    const {
        sorting,
        setSorting,
        columnFilters,
        setColumnFilters,
        columnVisibility,
        setColumnVisibility,
        rowSelection,
        setRowSelection,
        pagination: tablePagination,
        setPagination: setTablePagination,
        tableState,
    } = useTableState<TData>(data)

    React.useEffect(() => {
        if (!pagination) return

        setTablePagination(old => {
            if (old.pageIndex === pagination.page - 1 && old.pageSize === pagination.limit) {
                return old
            }
            return {pageIndex: pagination.page - 1, pageSize: pagination.limit}
        })
    }, [pagination])

    const sensors = useSensors(useSensor(MouseSensor), useSensor(TouchSensor), useSensor(KeyboardSensor))

    const dataIds = React.useMemo<UniqueIdentifier[]>(() => data.map(row => row.id), [data])

    const table = useReactTable({
        data,
        columns,
        state: {
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
            pagination: tablePagination,
        },
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        onColumnVisibilityChange: setColumnVisibility,
        onRowSelectionChange: setRowSelection,
        onPaginationChange: setTablePagination,
        manualPagination: true,
        pageCount: pagination ? Math.ceil(pagination.total / tableState.pagination.pageSize) : -1,
        getCoreRowModel: getCoreRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFacetedRowModel: getFacetedRowModel(),
        getFacetedUniqueValues: getFacetedUniqueValues(),
        getRowId: row => row.id,
        enableRowSelection: true,
    })

    const visibleColumns = table.getVisibleFlatColumns()
    const lastVisibleColumn = visibleColumns[visibleColumns.length - 1]
    const addClassForActionsLast = lastVisibleColumn?.id === "actions"

    function handleDragEnd(event: DragEndEvent) {
        const {active, over} = event
        if (active.id !== over?.id) {
            const oldIndex = dataIds.indexOf(active.id)
            const newIndex = dataIds.indexOf(over!.id)
            const newData = arrayMove(data, oldIndex, newIndex)
            setData(newData)
        }
    }

    return (
        <Tabs defaultValue="outline" className="w-full flex-col justify-start gap-6 py-8">
            <div className="flex items-center justify-between px-4 lg:px-6">
                <div className="flex items-center gap-2">
                    {customButtons?.map((btn, i) => (
                        <React.Fragment key={i}>{btn}</React.Fragment>
                    ))}
                </div>

                <div className="flex items-center gap-2">
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="outline" size="sm">
                                <IconLayoutColumns/>
                                <span className="hidden lg:inline">Columns</span>
                                <span className="lg:hidden">Columns</span>
                                <IconChevronDown/>
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-56">
                            {table
                                .getAllColumns()
                                .filter(col => col.getCanHide())
                                .map(column => (
                                    <DropdownMenuCheckboxItem
                                        key={column.id}
                                        className="capitalize"
                                        checked={column.getIsVisible()}
                                        onCheckedChange={value => column.toggleVisibility(!!value)}
                                    >
                                        {column.id}
                                    </DropdownMenuCheckboxItem>
                                ))}
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </div>

            <TabsContent value="outline" className="relative flex flex-col gap-4 overflow-auto px-4 lg:px-6">
                <div className="overflow-hidden rounded-lg border">
                    <DndContext
                        collisionDetection={closestCenter}
                        modifiers={[restrictToVerticalAxis]}
                        onDragEnd={handleDragEnd}
                        sensors={sensors}
                    >
                        <Table>
                            <TableHeader className="bg-muted sticky top-0 z-10">
                                {table.getHeaderGroups().map(headerGroup => (
                                    <TableRow key={headerGroup.id}>
                                        {headerGroup.headers.map((header, index) => (
                                            <TableHead key={header.id} colSpan={header.colSpan}
                                                       className={index === 0 ? "pl-5" : undefined}>
                                                {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                                            </TableHead>
                                        ))}
                                    </TableRow>
                                ))}
                            </TableHeader>
                            <TableBody
                                className={addClassForActionsLast ? "**:data-[slot=table-cell]:last:w-3 **:data-[slot=table-cell]:first:pl-5" : "**:data-[slot=table-cell]:first:pl-5"}
                            >
                                {table.getRowModel().rows.length ? (
                                    <SortableContext items={dataIds} strategy={verticalListSortingStrategy}>
                                        {table.getRowModel().rows.map(row => (
                                            <TableRow key={row.id} data-state={row.getIsSelected() && "selected"}>
                                                {row.getVisibleCells().map(cell => (
                                                    <TableCell key={cell.id}
                                                               className={cell.column.id === "actions" ? "actions" : undefined}>
                                                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                                    </TableCell>
                                                ))}
                                            </TableRow>
                                        ))}
                                    </SortableContext>
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={columns.length} className="h-24 text-center">
                                            No results.
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </DndContext>
                </div>

                <div className="flex items-center justify-between">
                    <div className="text-muted-foreground hidden flex-1 text-sm lg:flex">
                        Total count: {pagination?.total ?? data.length}
                    </div>
                    <div className="flex w-full items-center gap-8 lg:w-fit">
                        <div className="hidden items-center gap-2 lg:flex">
                            <Label className="text-sm font-medium">Rows per page</Label>
                            <Select
                                value={`${table.getState().pagination.pageSize}`}
                                onValueChange={value => {
                                    table.setPageSize(Number(value));
                                    onPaginationChange(table.getState().pagination.pageIndex, Number(value))
                                }}
                            >
                                <SelectTrigger size="sm" className="w-20">
                                    <SelectValue placeholder={`${table.getState().pagination.pageSize}`}/>
                                </SelectTrigger>
                                <SelectContent side="top">
                                    {[10, 20, 30, 40, 50].map(size => (
                                        <SelectItem key={size} value={`${size}`}>
                                            {size}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="flex w-fit items-center justify-center text-sm font-medium">
                            Page {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
                        </div>
                        <div className="ml-auto flex items-center gap-2 lg:ml-0">
                            <Button variant="outline" className="hidden h-8 w-8 p-0 lg:flex"
                                    onClick={() => {
                                        table.setPageIndex(0);
                                        onPaginationChange(0, table.getState().pagination.pageSize)
                                    }} disabled={!table.getCanPreviousPage()}>
                                <IconChevronsLeft/>
                            </Button>
                            <Button variant="outline" className="size-8" onClick={() => {
                                table.previousPage();
                                onPaginationChange(table.getState().pagination.pageIndex - 2, table.getState().pagination.pageSize)
                            }} disabled={!table.getCanPreviousPage()}>
                                <IconChevronLeft/>
                            </Button>
                            <Button variant="outline" className="size-8" onClick={() => {
                                table.nextPage();
                                onPaginationChange(table.getState().pagination.pageIndex + 2, table.getState().pagination.pageSize)
                            }} disabled={!table.getCanNextPage()}>
                                <IconChevronRight/>
                            </Button>
                            <Button variant="outline" className="hidden size-8 lg:flex"
                                    onClick={() => {
                                        table.setPageIndex(table.getPageCount() - 1);
                                        onPaginationChange(table.getPageCount(), table.getState().pagination.pageSize)
                                    }} disabled={!table.getCanNextPage()}>
                                <IconChevronsRight/>
                            </Button>
                        </div>
                    </div>
                </div>
            </TabsContent>
        </Tabs>
    )
}
