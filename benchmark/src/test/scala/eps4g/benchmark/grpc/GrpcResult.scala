package eps4g.benchmark.grpc

import io.gatling.core.Predef.Status

case class GrpcResult(status: Status, code: Option[String], message: Option[String])
