export interface Todo {
    id: number;
    title: string;
    description: string;
    status: 'pending' | 'in_progress' | 'completed';
    created_at: string;
    updated_at: string;
}

export interface TodoListResponse {
    data: Todo[];
    meta: {
        current_page: number;
        page_size: number;
        total_items: number;
        total_pages: number;
    };
}
