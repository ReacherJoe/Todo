import { useEffect, useState } from "react"
import "./styles.css"
import { NewTodoForm }  from "./NewTodoForm"
import TodoList from "./TodoList"
import axios from "axios"


export default function App () {
  const [todos, setTodos] = useState([]);
  


useEffect (() => {
  axios.get('http://localhost:8080/api/data')

  .then(response => setTodos(response.data))
  .catch(error => console.error('Error fetching data:', error));
}, []);


function addTodo(title) {
  const newTodo = { title, completed: false };

 
  axios.post('http://localhost:8080/api/data', newTodo)
  .then(response => {
    console.log('Response Data:', response.data);
    setTodos(currentTodos => [...currentTodos, response.data]);
  })
  .catch(error => {
    if (error.response.status === 400) {
      console.error('Validation Error:', error.response.data);
      // Handle validation error, show user-friendly message
    } else {
      console.error('Error adding todo:', error.response.data);
      // Handle other errors
    }
  });
  
  axios.interceptors.response.use(
    response => response,
    error => {
      console.error('Global Response Error:', error);
      return Promise.reject(error);
    }
  );

}


function toggleTodo(id, completed) {
  console.log(completed)
  setTodos(currentTodos => {
    const updatedTodos = currentTodos.map(todo => {
      if (todo.id === id) {
        return { ...todo, completed };
      }

      return todo;
    })
    console.log('updated Todos:', updatedTodos);

    return updatedTodos;
  })
  axios.put(`http://localhost:8080/api/data/${id}`, { completed })
  .then(response => {
  
    console.log('Todo updated successfully:', response.data);
  })
  .catch(error => {
   
    console.error('Error updating todo:', error);
    alert('An error occurred while updating the todo. Reverting changes.');
   
    fetchTodos(); 
  });
}



function deleteTodo(id) {
 


   axios
      .delete(`http://localhost:8080/api/data/${id}`)
      .then(() => {
        // If the delete request is successful, update the state
        setTodos((currentTodos) =>
          currentTodos.filter((todo) => todo.id !== id)
        );
      })
      .catch((error) => console.error("Error deleting todo:", error));
  
  
}

  return (
  <>
  <NewTodoForm onSubmit={addTodo}/>
  <h1 className="header">Todo List</h1>
  <TodoList todos={todos} toggleTodo={toggleTodo} deleteTodo={deleteTodo}/>
  </>
  )
}
