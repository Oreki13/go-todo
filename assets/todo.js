const saveTodo = event => {
  const data = {
    title: document.getElementById("input-todo").value,
    completed: 0
  };
  const url = "http://localhost:8080/todo";
  const val = document.getElementById("input-todo").value

  if (event.keyCode === 13 && val.length > 1) {
    $.post(url, data, function (data, status) {
      console.log(`${data} and status is ${status}`);
      location.reload();
    });

  }
};

const goCheck = (e) => {
  $.ajax({
    type: `PATCH`,
    url: `http://localhost:8080/todos/${e}`,
    processData: false,
    contentType: false,
  })
};

