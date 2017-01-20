
CREATE EXTENSION "uuid-ossp";


create table if not exists t_wc_user(
	id varchar(36) primary key default(uuid_generate_v4()),
	userid varchar(100)  unique not null,	--成员UserID。对应管理端的帐号
	name varchar(100),	--成员名称
	department varchar(100),	--成员所属部门id列表
	mobile	varchar(20),	--手机号码
	gender	varchar(1),	--性别。0表示未定义，1表示男性，2表示女性
	email		varchar(100),	--邮箱
	weixinid varchar(100),	--微信号
	avatar	varchar(300),	--头像url。注：如果要获取小图将url最后的"/0"改成"/64"即可
	status	varchar(1),		--关注状态: 1=已关注，2=已冻结，4=未关注
	product varchar(100)[],	--关注的产品code
	crt timestamp,	--创建时间
	lut timestamp,	--更新时间
	del boolean default('false')	--是否删除
);

create table if not exists t_wc_product(
	id varchar(36) primary key default(uuid_generate_v4()),
	productcode varchar(100), --产品编码
	crt timestamp,	--创建时间
	lut timestamp,	--更新时间
	del boolean default('false')	--是否删除
);

create table if not exists t_wc_news(
	id varchar(36) primary key default(uuid_generate_v4()),
	title varchar(1000),					--标题
	summary varchar(1000),				--摘要
	content text,									--内容
	productcode varchar(100),			--产品编码
	crt timestamp DEFAULT now(),								--创建时间
	del boolean default('false')	--是否删除
);

create table if not exists t_wc_admin(
	id varchar(36) primary key default(uuid_generate_v4()),
	userid varchar(100)  unique not null,	--管理端的帐号
	passwd varchar(100), --管理员密码
	name varchar(100),	--成员名称
	crt timestamp,	--创建时间
	lut timestamp,	--更新时间
	del boolean default('false')	--是否删除
);

alter table t_wc_product
add brief varchar(1000);
comment on column t_wc_product.brief is '产品简介';

--20151027
--publish news by tag
ALTER TABLE t_wc_product RENAME TO t_wc_tag;
alter table t_wc_tag rename column productcode to tagcode;
comment on column t_wc_tag.tagcode is '标签编码';
alter table t_wc_user rename column product to tags;
comment on column t_wc_user.tags is '标签';
alter table t_wc_news rename column productcode to tagcode;
comment on column t_wc_news.tagcode is '标签编码';







