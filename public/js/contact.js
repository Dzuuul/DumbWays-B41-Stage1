function redirectToEmail() {
  let getName = document.getElementById("input-name").value;
  let getEmail = document.getElementById("input-email").value;
  let getPhoneNumber = document.getElementById("input-phonenumber").value;
  let getSubject = document.getElementById("input-subject").value;
  let getMessage = document.getElementById("input-yourmessage").value;

  if (getName == "") {
    return alert("Name Required!");
  } else if (getEmail == "") {
    return alert("Email Required!");
  } else if (getPhoneNumber == "") {
    return alert("Phone Number is Required!");
  } else if (getSubject == "") {
    return alert("Choose the subject!");
  } else if (getMessage == "") {
    return alert("Your Message is Required!");
  }

  let mailTo = "rhomairama.dev@gmail.com";

  let a = document.createElement("a");

  a.href = `mailto:${mailTo}?subject=${getSubject}&body=Hello my name ${getName},%0D%0A${getMessage}.%0D%0APlease contact me ASAP at ${getPhoneNumber}.%0D%0AThank You!%0D%0A%0D%0A%0D%0A`;
  a.click();
}
