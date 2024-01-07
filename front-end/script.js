function submitForm() {
    const fullName = document.getElementById('fullName').value;
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const phoneNumber = document.getElementById('phoneNumber').value;
    const password = document.getElementById('password').value;
    const gender = document.querySelector('input[name="gender"]:checked').id;

    const userData = {
      full_name: fullName,
      username: username,
      email: email,
      phone_number: phoneNumber,
      password: password,
      gender: gender
    };

    fetch('http://localhost:5050/uploadjson', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      })
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        console.log('Success:', data);
        alert('Registration successful!');
      })
      .catch(error => {
        console.error('Error:', error);
        const errorBox = document.createElement('div');
        errorBox.textContent = 'Error registering user: ' + error.message;
        errorBox.style = 'position: fixed; top: 10px; left: 10px; padding: 10px; background-color: #f00; color: #fff;';
        document.body.appendChild(errorBox);
        setTimeout(() => {
          document.body.removeChild(errorBox);
        }, 5000);
      });
  }

  function getUsersAndDownloadCSV() {
    fetch('http://localhost:5050/getusers')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.blob();
      })
      .then(blob => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'users.csv');
        document.body.appendChild(link);
        link.click();
        link.remove();
      })
      .catch(error => {
        console.error('Error fetching users:', error);
        const errorBox = document.createElement('div');
        errorBox.textContent = 'Error fetching users: ' + error.message;
        errorBox.style = 'position: fixed; top: 10px; left: 10px; padding: 10px; background-color: #f00; color: #fff;';
        document.body.appendChild(errorBox);
        setTimeout(() => {
          document.body.removeChild(errorBox);
        }, 5000);
      });
  }