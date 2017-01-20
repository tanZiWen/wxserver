DROP TABLE IF EXISTS app_user;
CREATE TABLE app_user (
  id         BIGINT PRIMARY KEY,
  name       VARCHAR(100),
  mobile     VARCHAR(20),
  note       VARCHAR(3000),
  number     INTEGER,
  crt        TIMESTAMP WITH TIME ZONE,
  lut        TIMESTAMP WITH TIME ZONE,
  status     SMALLINT DEFAULT (1)
);

COMMENT ON TABLE app_user IS '用户表';
COMMENT ON COLUMN app_user.id IS '主键';
COMMENT ON COLUMN app_user.name IS '用户姓名';
COMMENT ON COLUMN app_user.mobile IS '手机号码';
COMMENT ON COLUMN app_user.note IS '备注';
COMMENT ON COLUMN app_user.number IS '预约膏方人数';
COMMENT ON COLUMN app_user.crt IS '创建时间';
COMMENT ON COLUMN app_user.lut IS '最后更新时间';
COMMENT ON COLUMN app_user.status IS '状态，0:删除, 1:正常';

