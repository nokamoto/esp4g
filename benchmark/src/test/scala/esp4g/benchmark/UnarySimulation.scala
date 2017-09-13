package eps4g.benchmark

import eps4g.benchmark.grpc.{GrpcActionBuilder, GrpcProtocol}
import esp4g.benchmark.{Unary, UnaryServiceGrpc}
import io.gatling.core.Predef._
import io.grpc.netty.NettyChannelBuilder

import scala.concurrent.duration._
import scala.util.Try

abstract class UnarySimulation(size: Int, rps: Int) extends Simulation {
  private[this] val over: FiniteDuration = 60.seconds

  private[this] val protocol = new GrpcProtocol

  private[this] val channel = NettyChannelBuilder.forAddress("localhost", 9000).usePlaintext(true).build()

  private[this] val stub = UnaryServiceGrpc.stub(channel)

  // warm up
  (1 to 10).foreach(_ => Try(UnaryServiceGrpc.blockingStub(channel).send(Unary())))

  private[this] val unaryScenario = {
    scenario(s"Unary gRPC Proxy: $size bytes $rps rps").exec(new GrpcActionBuilder(f = { _ =>
      val unary = Unary().update(_.text := scala.util.Random.alphanumeric.take(size).mkString)
      stub.send(unary)
    })).inject(constantUsersPerSec(rps) during over)
  }

  setUp(unaryScenario).protocols(protocol)
}

class UnarySimulation128b100r extends UnarySimulation(128, 100)
class UnarySimulation128b1000r extends UnarySimulation(128, 1000)
