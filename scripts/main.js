var pickAction = "Pick your action!"

function init() {
  swal({
      title: "Is it your turn?",
      text: "It's hellojessemillar@gmail.com's turn. Are you hellojessemillar@gmail.com?",
      icon: "warning",
      buttons: ["Nope!", "I'm the droid you're looking for."],
    })
    .then((willDelete) => {
      if (willDelete) {
        swal(pickAction, {
          icon: "success",
        });
      } else {
        swal("Please wait for your turn. You'll receive another email when it's time.").then(() => {
          wipePage();
        });
      }
    });
}

function stand() {
  swal({
      title: "You want to stand?",
      text: "You're sure you want to do nothing for your turn?",
      icon: "warning",
      buttons: ["Nah.", "Yeah!"],
    })
    .then((willDelete) => {
      if (willDelete) {
        swal("You chose to stand.", {
          icon: "success",
        }).then(() => {
          turnOver()
        });
      } else {
        swal(pickAction);
      }
    });
}

function trash() {
  swal({
      title: "Are you sure you want to trash?",
      text: "If you trash, you drop out of the game permanently.",
      icon: "warning",
      buttons: ["Nope!", "Yes."],
      dangerMode: true,
    })
    .then((willDelete) => {
      if (willDelete) {
        swal("You've withdrawn from the game.", {
          icon: "success",
          button: "'Til the Spire.",
        }).then(() => {
          wipePage()
        });
      } else {
        swal("You're still in the game!");
      }
    });
}

function turnOver() {
  swal({
    title: "Turn over.",
    text: "Your turn is now over! Please wait for the next email.",
    icon: "success",
    button: "Patience, young padawan.",
  }).then(
    wipePage()
  );
}

function wipePage() {
  document.body.innerHTML = '';
  window.close();
}
