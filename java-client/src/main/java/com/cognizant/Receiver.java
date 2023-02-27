package com.cognizant;

import java.io.IOException;
import java.io.InputStream;
import java.net.Socket;
import java.nio.charset.StandardCharsets;

public class Receiver extends Thread{

    private final String url;

    private final int port;

    public Receiver(String url, int port){
        this.url = url;
        this.port = port;
    }

    @Override
    public synchronized void start() {
        while (true) {
            try (Socket socket = new Socket(this.url, this.port)) {
                InputStream in = socket.getInputStream();
                System.out.println("reading started ...");
                String login = "login:!javaclient!";
                socket.getOutputStream().write(login.getBytes(StandardCharsets.UTF_8));
                byte[] buffer = new byte[1024];
                int bytes;
                while((bytes = in.read(buffer)) > 0){
                    System.out.print(bytes + " bytes message: ");
                    byte[] received = new byte[bytes];
                    System.arraycopy(buffer, 0, received, 0, bytes);
                    System.out.println(new String(received, StandardCharsets.UTF_8));
                }
            } catch (IOException e) {
                System.out.println(e.getMessage());
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException ignored) {

                }
            }
        }
    }
}
