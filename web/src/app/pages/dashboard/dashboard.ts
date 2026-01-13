import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Observable } from 'rxjs';
import { TodoService } from '../../services/todo';
import { AuthService } from '../../services/auth';
import { Todo } from '../../models/todo.model';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
})
export class Dashboard implements OnInit {
  todos = signal<Todo[]>([]);
  isLoading = signal<boolean>(false);
  showAddForm = signal<boolean>(false);

  todoForm: FormGroup;
  currentUser: Observable<User | null>;

  constructor(
    private todoService: TodoService,
    private authService: AuthService,
    private fb: FormBuilder
  ) {
    this.currentUser = this.authService.currentUser$;
    this.todoForm = this.fb.group({
      title: ['', [Validators.required, Validators.minLength(3)]],
      description: [''],
      due_date: ['']
    });
  }

  ngOnInit() {
    this.loadTodos();
  }

  loadTodos() {
    this.isLoading.set(true);
    this.todoService.getAll().subscribe({
      next: (res) => {
        if (res.success && res.data) {
          // Flatten the response struct handled by HttpClient generic but API returns {data: [...]}
          // Wait, my ApiResponse<TodoListResponse> means res.data is TodoListResponse
          // And TodoListResponse has .data which is Todo[]
          // So it's res.data.data
          this.todos.set(res.data.data);
        }
        this.isLoading.set(false);
      },
      error: () => this.isLoading.set(false)
    });
  }

  addTodo() {
    if (this.todoForm.invalid) return;

    this.todoService.create(this.todoForm.value).subscribe({
      next: (res) => {
        if (res.success && res.data) {
          this.todos.update(todos => [res.data!, ...todos]);
          this.todoForm.reset();
          this.showAddForm.set(false);
        }
      }
    });
  }

  toggleStatus(todo: Todo) {
    const newStatus = todo.status === 'completed' ? 'pending' : 'completed';
    this.todoService.update(todo.id, { status: newStatus }).subscribe({
      next: (res) => {
        if (res.success && res.data) {
          this.todos.update(todos =>
            todos.map(t => t.id === todo.id ? res.data! : t)
          );
        }
      }
    });
  }

  deleteTodo(id: number) {
    if (!confirm('Are you sure?')) return;

    this.todoService.delete(id).subscribe({
      next: () => {
        this.todos.update(todos => todos.filter(t => t.id !== id));
      }
    });
  }

  logout() {
    this.authService.logout();
  }
}

