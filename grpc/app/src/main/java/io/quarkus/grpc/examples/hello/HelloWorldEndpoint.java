package io.quarkus.grpc.examples.hello;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;

import org.eclipse.microprofile.config.inject.ConfigProperty;

import examples.Greeter;
import examples.GreeterGrpc;
import examples.HelloReply;
import examples.HelloRequest;
import io.quarkus.grpc.GrpcClient;
import io.smallrye.mutiny.Uni;

@Path("/hello")
public class HelloWorldEndpoint {

    @GrpcClient("hello")
    GreeterGrpc.GreeterBlockingStub blockingHelloService;

    @GrpcClient("hello")
    Greeter helloService;

    @ConfigProperty(name = "numrequest")
    int numRequest;

    @ConfigProperty(name = "teststring")
    String testString;

    @ConfigProperty(name = "quarkus.grpc.clients.hello.host")
    String host;

    @GET
    @Path("/blocking/{name}")
    public String helloBlocking(String name) {
        HelloReply reply = blockingHelloService.sayHello(HelloRequest.newBuilder().setName(name).build());
        return generateResponse(reply);

    }

    @GET
    @Path("/mutiny/{name}")
    public Uni<String> helloMutiny(String name) {
        return helloService.sayHello(HelloRequest.newBuilder().setName(name).build())
                .onItem().transform((reply) -> generateResponse(reply));
    }

    public String generateResponse(HelloReply reply) {
        return String.format("%s! HelloWorldService has been called %d number of times.", reply.getMessage(), reply.getCount());
    }

    @GET
    @Path("/grpc")
    public String testGRPC() throws InterruptedException
    {
        MyClient client = new MyClient(host, 9000);
        try {
            client.makeMultipleRequests(testString, numRequest);
        } finally {
            client.shutdown();
        }
        return "Hello";
    }
}

