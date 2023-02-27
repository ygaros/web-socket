package com.cognizant;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello world!");
        Receiver receiver = new Receiver("localhost", 9876);
        receiver.start();
    }
}