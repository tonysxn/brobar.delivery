export type AuthData = {
    access: {
        expires_at: string;
        token: string;
    };
    refresh: {
        expires_at: string;
        token: string;
    };
    user: {
        id: string;
        role_id: string;
        email: string;
        name: string;
        address: string;
        address_coords: string;
        phone: string;
        promo_card: string;
    };
};
