{ pkgs ? import <nixpkgs> {} }:

pkgs.stdenv.mkDerivation rec {
	name = "gotk-layer-shell";
	version = "0.0.1";

	buildInputs = with pkgs; [
		gnome3.glib
		gnome3.gtk
		gtk-layer-shell.dev
	];

	nativeBuildInputs = with pkgs; [
		pkgconfig go
	];
}
