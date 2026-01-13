export interface User {
    id: number;
    name: string;
    email: string;
    created_at: string;
}

export interface AuthResponse {
    user: User;
    tokens: {
        access_token: string;
        refresh_token: string;
    };
}
