let dataProject = [];

function addProject(event) {
  event.preventDefault();
  let projectName = document.getElementById("input-project-name").value;
  let duration = dateDuration();
  let description = document.getElementById("input-description").value;

  let reactJs = document.getElementById("input-reactJs").checked;
  let vueJs = document.getElementById("input-vueJs").checked;
  let angular = document.getElementById("input-angular").checked;
  let laravel = document.getElementById("input-laravel").checked;

  let uploadImage = document.getElementById("input-upload-img").files;

  if (reactJs) {
    reactJs = '<i class="fa-brands fa-react fa-xl"></i>';
  } else {
    reactJs = "";
  }
  if (vueJs) {
    vueJs = '<i class="fa-brands fa-vuejs fa-xl"></i>';
  } else {
    vueJs = "";
  }
  if (angular) {
    angular = '<i class="fa-brands fa-angular fa-xl"></i>';
  } else {
    angular = "";
  }
  if (laravel) {
    laravel = '<i class="fa-brands fa-laravel fa-xl"></i>';
  } else {
    laravel = "";
  }

  if (projectName == "") {
    return alert("Please fill the project name 🙃");
  } else if (typeof duration === "undefined") {
    return alert("Please enter the date correctly 🙏🏻");
  } else if (description == "") {
    return alert("Description cannot be empty 🙌🏻");
  } else if (reactJs == "" && vueJs == "" && angular == "" && laravel == "") {
    return alert("At least check one technologies ✅⚙️");
  }

  if (uploadImage.length != 0) {
    uploadImage = URL.createObjectURL(uploadImage[0]);
  } else {
    return alert("Upload your picture don't be shy 📸");
  }

  let projectItem = {
    projectName,
    duration,
    description,
    reactJs,
    vueJs,
    angular,
    laravel,
    uploadImage,
  };

  dataProject.push(projectItem);
  console.log(dataProject);

  renderProject();
}

function renderProject() {
  document.getElementById("project").innerHTML = ``;

  for (let i = 0; i < dataProject.length; i++) {
    document.getElementById("project").innerHTML += `
    
  <div class="project-item">
  <div class="project-image">
    <img src="${dataProject[i].uploadImage}"
      alt="project image">
  </div>
  <a href="/pages/detail-project.html">
    <h3>${dataProject[i].projectName}</h3>
  </a>
  <p class="duration-p">duration : ${dataProject[i].duration}</p>
  <div class="item-p-container">
    <p>${dataProject[i].description}</p>
  </div>
  <div class="item-icons">
    ${dataProject[i].reactJs}
    ${dataProject[i].vueJs}
    ${dataProject[i].angular}
    ${dataProject[i].laravel}

  </div>
  <div class="button-group">
    <button>edit</button>
    <button>delete</button>
  </div>
</div>
  `;
  }
}

function dateDuration() {
  let sD = document.getElementById("input-start-date").value;
  let eD = document.getElementById("input-end-date").value;
  let startDate = new Date(sD);
  let endDate = new Date(eD);

  let timeDifference = endDate.getTime() - startDate.getTime();
  let monthDifference = Math.floor(
    timeDifference / (1000 * 3600 * 24 * 30.4167)
  );
  let dateDifference = Math.floor(timeDifference / (1000 * 3600 * 24));

  if (dateDifference == 1 || dateDifference == 0) {
    return `${dateDifference} day`;
  } else if (dateDifference > 0 && dateDifference < 30.4167) {
    return `${dateDifference} days`;
  } else if (dateDifference == 30.4167) {
    return `${monthDifference} month`;
  } else if (dateDifference > 30.4167) {
    return `${monthDifference} months`;
  }
}
