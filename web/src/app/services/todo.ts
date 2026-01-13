import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Todo, TodoListResponse } from '../models/todo.model';
import { ApiResponse } from '../models/api-response.model';

@Injectable({
  providedIn: 'root',
})
export class TodoService {
  private apiUrl = 'http://localhost:8080/api/todos';

  constructor(private http: HttpClient) { }

  getAll(page = 1, pageSize = 10, status = ''): Observable<ApiResponse<TodoListResponse>> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    if (status) {
      params = params.set('status', status);
    }

    return this.http.get<ApiResponse<TodoListResponse>>(this.apiUrl, { params });
  }

  create(todo: Partial<Todo>): Observable<ApiResponse<Todo>> {
    return this.http.post<ApiResponse<Todo>>(this.apiUrl, todo);
  }

  update(id: number, todo: Partial<Todo>): Observable<ApiResponse<Todo>> {
    return this.http.put<ApiResponse<Todo>>(`${this.apiUrl}/${id}`, todo);
  }

  delete(id: number): Observable<ApiResponse<void>> {
    return this.http.delete<ApiResponse<void>>(`${this.apiUrl}/${id}`);
  }
}

