package eps4g.benchmark.grpc

import com.trueaccord.scalapb.GeneratedMessage
import io.gatling.core.action.Action
import io.gatling.core.action.builder.ActionBuilder
import io.gatling.core.session.Session
import io.gatling.core.structure.ScenarioContext

import scala.concurrent.Future

class GrpcActionBuilder(f: Session => Future[GeneratedMessage]) extends ActionBuilder {
  override def build(ctx: ScenarioContext, next: Action): Action = {
    new GrpcAction(f = f, stats = ctx.coreComponents.statsEngine, next = next, system = ctx.system)
  }
}
