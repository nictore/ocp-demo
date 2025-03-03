package io.quarkus.grpc.examples.hello;

import examples.HelloReply;
import examples.HelloRequest;
import io.grpc.ManagedChannel;
import examples.GreeterGrpc;
import io.grpc.ManagedChannelBuilder;
import java.util.concurrent.TimeUnit;

public class MyClient {
    private final ManagedChannel channel;
    private final GreeterGrpc.GreeterBlockingStub blockingStub;

    public MyClient(String host, int port) {
        this(ManagedChannelBuilder.forAddress(host, port).usePlaintext().build());
    }

    private MyClient(ManagedChannel channel) {
        this.channel = channel;
        this.blockingStub = GreeterGrpc.newBlockingStub(channel);
    }

    public void shutdown() throws InterruptedException {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
    }

    public void makeMultipleRequests(String name, int numRequests) {
        for (int i = 0; i < numRequests; i++) {
            HelloRequest request = HelloRequest.newBuilder().setName(name + i).build();
            HelloReply response = blockingStub.sayHello(request);
            System.out.println(response.getMessage());
        }
    }
}