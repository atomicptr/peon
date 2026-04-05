{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    tparse
  ];

  PEON_DEBUG = "true";
}
