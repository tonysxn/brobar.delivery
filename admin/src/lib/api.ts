import axios from 'axios';

const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL || 'http://localhost:8000';

export interface Setting {
    key: string;
    setting_type: string;
    value: string;
}

export interface Review {
    id: string;
    food_rating: number;
    service_rating: number;
    comment: string;
    phone?: string;
    email?: string;
    name?: string;
    created_at: string;
}

export const settingsApi = {
    getAll: async () => {
        const { data } = await axios.get<{ success: boolean; data: Setting[] }>(`${GATEWAY_URL}/settings`);
        return data.data;
    },
    update: async (key: string, value: string, type: string, token?: string) => {
        const headers = token ? { Authorization: `Bearer ${token}` } : {};
        const { data } = await axios.put(`${GATEWAY_URL}/settings/${key}`, { value, type }, { headers });
        return data;
    },
};

export const reviewsApi = {
    getAll: async () => {
        const { data } = await axios.get<{ success: boolean; data: Review[] }>(`${GATEWAY_URL}/reviews`);
        return data.data;
    },
    delete: async (id: string, token?: string) => {
        const headers = token ? { Authorization: `Bearer ${token}` } : {};
        const { data } = await axios.delete(`${GATEWAY_URL}/reviews/${id}`, { headers });
        return data;
    },
};

export interface Order {
    id: string;
    name: string;
    address: string;
    phone: string;
    status_id: string;
    total_price: number;
    created_at: string;
    delivery_type_id: string;
    payment_method: string;
    zone?: string;
    delivery_cost: number;
    delivery_door: boolean;
    delivery_door_price?: number;
    items?: OrderItem[];
    wishes?: string;
}

export interface OrderItem {
    id: string;
    product_id: string;
    external_product_id: string;
    name: string;
    price: number;
    quantity: number;
    total_price: number;
    weight: number;
    total_weight: number;
    product_variation_group_name?: string;
    product_variation_name?: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    data: T[];
    pagination: {
        total_count: number;
        page: number;
        limit: number;
        order_by: string;
        order_dir: string;
    };
}

export const ordersApi = {
    getAll: async (page = 1, limit = 20, token?: string) => {
        const headers = token ? { Authorization: `Bearer ${token}` } : {};
        const { data } = await axios.get<PaginatedResponse<Order>>(
            `${GATEWAY_URL}/orders?page=${page}&limit=${limit}`,
            { headers }
        );
        return data;
    },
    getOne: async (id: string, token?: string) => {
        const headers = token ? { Authorization: `Bearer ${token}` } : {};
        const { data } = await axios.get<{ success: boolean; data: Order }>(
            `${GATEWAY_URL}/orders/${id}`,
            { headers }
        );
        return data.data;
    },
};
