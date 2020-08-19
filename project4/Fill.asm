// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

@prev
M=0

(LOOP)
  @KBD
  D=M

  @WHITE
  D;JEQ
  @BLACK
  D;JNE
  
(WHITE)
  @color
  M=0
  @prev
  D=M
  @FILL
  D;JNE
  @LOOP
  0;JMP

(BLACK)
  @color
  M=-1
  @prev
  D=M
  @LOOP
  D;JNE

(FILL)
  @color
  D=M
  @prev
  M=D
  @KBD
  D=A
  @n
  M=D
  @SCREEN
  D=A
  @addr
  M=D
  (FILLLOOP)
    @n
    D=M
    @addr
    D=D-M
    @LOOP
    D;JEQ

    @color
    D=M

    @addr
    A=M
    M=D

  	@addr
  	M=M+1

    @FILLLOOP
    0;JMP
