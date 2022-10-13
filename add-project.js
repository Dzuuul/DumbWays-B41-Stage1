let dataProject = [];

function addProject(event) {
  event.preventDefault();
  let projectName = document.getElementById("input-project-name").value;
  let description = document.getElementById("input-description").value;
  let uploadImage = document.getElementById("input-upload-img").files;

  let reactJs = document.getElementById("input-reactJs").checked;
  let vueJs = document.getElementById("input-vueJs").checked;
  let angular = document.getElementById("input-angular").checked;
  let laravel = document.getElementById("input-laravel").checked;

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

  uploadImage = URL.createObjectURL(uploadImage[0]);

  let projectItem = {
    projectName,
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
  <a href="/pages/detail-project.html" target="_blank">
    <h3>${dataProject[i].projectName}</h3>
  </a>
  <p class="duration-p">durasi : -</p>
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
