package eps4g.benchmark.grpc

import akka.actor.ActorSystem
import com.trueaccord.scalapb.GeneratedMessage
import io.gatling.commons.stats.{KO, OK}
import io.gatling.core.action.Action
import io.gatling.core.session.Session
import io.gatling.core.stats.StatsEngine
import io.gatling.core.stats.message.ResponseTimings
import io.grpc.{StatusException, StatusRuntimeException}

import scala.concurrent.{ExecutionContextExecutor, Future}
import scala.util.{Failure, Success}

class GrpcAction(f: Session => Future[GeneratedMessage], stats: StatsEngine, next: Action, system: ActorSystem) extends Action {
  override def name: String = "grpc"

  override def execute(session: Session): Unit = {
    implicit val ctx: ExecutionContextExecutor = system.dispatcher

    val requestName = ""

    val start = System.currentTimeMillis()

    f(session)
      .map(_ => GrpcResult(OK, None, None))
      .recover {
        case e: StatusException => GrpcResult(KO, Some(e.getStatus.toString), Some(e.getMessage))
        case e: StatusRuntimeException => GrpcResult(KO, Some(e.getStatus.toString), Some(e.getMessage))
      }
      .onComplete {
        case Failure(e) =>
          stats.logCrash(session, requestName, e.getMessage)

        case Success(res) =>
          val end = System.currentTimeMillis()

          val timings = ResponseTimings(startTimestamp = start, endTimestamp = end)

          stats.logResponse(
            session = session,
            requestName = requestName,
            status = res.status,
            timings = timings,
            responseCode = res.code,
            message = res.message
          )
      }

    f(session).onComplete(_ => next ! session)
  }
}
