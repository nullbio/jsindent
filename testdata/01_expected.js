this.calculatorForm = function() {
  return DIV({class: "row"},
    DIV({class: "col-lg-12"},
      DIV({class: "ibox float-e-margins"}, [
        DIV({class: "ibox-title"}, [
          H5("New Profile ", m("small", "Create a new calculator Betting Profile")),
          iboxTools
        ]),
        DIV({class: "ibox-content"},
          DIV({class: "form-horizontal"})
        )
      ]);
    )
  )
};
